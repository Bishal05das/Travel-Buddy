import { useId, type InputHTMLAttributes, type ReactNode, type SelectHTMLAttributes, type TextareaHTMLAttributes } from "react";

const control =
  "w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:border-teal-600 focus:outline-none focus:ring-2 focus:ring-teal-600/20 disabled:bg-slate-100";

function Wrapper({ id, label, hint, error, children }: { id: string; label: string; hint?: ReactNode; error?: string; children: ReactNode }) {
  return (
    <div className="space-y-1">
      <label htmlFor={id} className="block text-sm font-medium text-slate-700">
        {label}
      </label>
      {children}
      {error ? (
        <p id={`${id}-error`} className="text-xs text-red-600">
          {error}
        </p>
      ) : hint ? (
        <p className="text-xs text-slate-500">{hint}</p>
      ) : null}
    </div>
  );
}

interface Common {
  label: string;
  hint?: ReactNode;
  error?: string;
}

export function Input({ label, hint, error, ...rest }: Common & InputHTMLAttributes<HTMLInputElement>) {
  const id = useId();
  return (
    <Wrapper id={id} label={label} hint={hint} error={error}>
      <input id={id} className={control} aria-invalid={!!error} aria-describedby={error ? `${id}-error` : undefined} {...rest} />
    </Wrapper>
  );
}

export function TextArea({ label, hint, error, ...rest }: Common & TextareaHTMLAttributes<HTMLTextAreaElement>) {
  const id = useId();
  return (
    <Wrapper id={id} label={label} hint={hint} error={error}>
      <textarea id={id} className={control} rows={4} aria-invalid={!!error} {...rest} />
    </Wrapper>
  );
}

export function Select({ label, hint, error, children, ...rest }: Common & SelectHTMLAttributes<HTMLSelectElement>) {
  const id = useId();
  return (
    <Wrapper id={id} label={label} hint={hint} error={error}>
      <select id={id} className={control} aria-invalid={!!error} {...rest}>
        {children}
      </select>
    </Wrapper>
  );
}
