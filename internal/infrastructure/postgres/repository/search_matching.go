package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
)

type searchColumn struct {
	name   string
	weight int
}

func escapeSearchPattern(text string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(text)
}

func fuzzySearchTerm(term string) bool {
	// These characters are literal user input, not LIKE wildcards. Trigrams
	// ignore punctuation, so fuzzy matching would lose that distinction.
	if strings.ContainsAny(term, `%_\`) {
		return false
	}
	letters := 0
	for _, r := range term {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			letters++
		}
	}
	return letters >= 3
}

// Every search word must match one field. This permits reordered words and
// combinations such as destination + agency while keeping filters conjunctive.
func (s *searchSQL) matches(query string, columns []searchColumn) []string {
	conditions := []string{}
	for _, term := range strings.Fields(query) {
		pattern := s.bind("%" + escapeSearchPattern(term) + "%")
		parts := []string{}
		for _, column := range columns {
			parts = append(parts, fmt.Sprintf(`%s ILIKE %s ESCAPE E'\\'`, column.name, pattern))
		}
		if fuzzySearchTerm(term) {
			word := s.bind(term)
			for _, column := range columns {
				parts = append(parts, word+" <% "+column.name)
			}
		}
		conditions = append(conditions, "("+strings.Join(parts, " OR ")+")")
	}
	return conditions
}

// Exact titles precede prefixes and substrings. Agency matches follow title
// matches, then weighted word similarity orders fuzzy and description matches.
func (s *searchSQL) ranking(query string, columns []searchColumn) string {
	if query == "" {
		return ""
	}
	exact := s.bind(escapeSearchPattern(query))
	prefix := s.bind(escapeSearchPattern(query) + "%")
	contains := s.bind("%" + escapeSearchPattern(query) + "%")
	classes := []string{}
	for i, column := range columns {
		if column.name == "t.description" {
			continue
		}
		for j, pattern := range []string{exact, prefix, contains} {
			classes = append(classes, fmt.Sprintf(`WHEN %s ILIKE %s ESCAPE E'\\' THEN %d`, column.name, pattern, i*3+j))
		}
	}
	scores := []string{}
	for _, term := range strings.Fields(query) {
		word := s.bind(term)
		fields := []string{}
		for _, column := range columns {
			fields = append(fields, fmt.Sprintf("%d * word_similarity(%s,%s)", column.weight, word, column.name))
		}
		if len(fields) == 1 {
			scores = append(scores, fields[0])
		} else {
			scores = append(scores, "GREATEST("+strings.Join(fields, ",")+")")
		}
	}
	return "CASE " + strings.Join(classes, " ") + " ELSE 9 END ASC,(" + strings.Join(scores, "+") + ") DESC,"
}

func (r *searchRepository) withSearch(ctx context.Context, query string, fn func(sqlx.ExtContext) error) error {
	fuzzy := false
	for _, term := range strings.Fields(query) {
		if fuzzySearchTerm(term) {
			fuzzy = true
			break
		}
	}
	if !fuzzy {
		return fn(r.db)
	}
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// SET LOCAL stays on this transaction's connection and cannot leak to
	// another request when the connection returns to the pool.
	if _, err := tx.ExecContext(ctx, `SET LOCAL pg_trgm.word_similarity_threshold = 0.45`); err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
