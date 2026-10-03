"use client";

import { useEffect, useRef, useState } from "react";
import { ApiError } from "@/lib/api";
import { imageUrl } from "@/lib/config";
import type { ImageUpdateResult } from "@/lib/endpoints";
import { Button } from "./ui/Button";
import { Alert } from "./ui/Feedback";
import { Input } from "./ui/Field";

export function ImageEditor({ imagePath, label, permission, onSave }: {
  imagePath?: string;
  label: string;
  permission: string;
  onSave: (form: FormData) => Promise<ImageUpdateResult>;
}) {
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [savedPath, setSavedPath] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [saving, setSaving] = useState(false);
  const formElement = useRef<HTMLFormElement>(null);

  useEffect(() => {
    if (!file) return;
    const url = URL.createObjectURL(file);
    // eslint-disable-next-line react-hooks/set-state-in-effect -- preview exists only while this file is selected
    setPreview(url);
    return () => URL.revokeObjectURL(url);
  }, [file]);

  function clearSelection() {
    setFile(null);
    setPreview(null);
    formElement.current?.reset();
  }

  function select(next: File | null) {
    setError(null);
    setSuccess(false);
    setFile(null);
    setPreview(null);
    if (!next) return;
    if (!/\.(jpe?g|png|webp)$/i.test(next.name) || !["image/jpeg", "image/png", "image/webp"].includes(next.type)) {
      setError("Choose a JPG, PNG or WebP image.");
      formElement.current?.reset();
    } else if (next.size === 0 || next.size > 10 * 1024 * 1024) {
      setError("Choose an image up to 10 MB. Empty files are not allowed.");
      formElement.current?.reset();
    } else {
      setFile(next);
    }
  }

  async function save(e: React.FormEvent) {
    e.preventDefault();
    if (!file || saving) return;
    setSaving(true);
    setError(null);
    setSuccess(false);
    const form = new FormData();
    form.set("image", file);
    try {
      const result = await onSave(form);
      setSavedPath(result.image_path);
      clearSelection();
      setSuccess(true);
    } catch (err) {
      setError(err instanceof ApiError && err.status === 403
        ? `Ask your agency owner to grant you the ${permission} permission to change this image.`
        : err instanceof Error ? err.message : "Could not save the image.");
    } finally {
      setSaving(false);
    }
  }

  const src = (file && preview) || imageUrl(savedPath ?? imagePath);
  return (
    <form ref={formElement} onSubmit={save} className="max-w-2xl space-y-4 rounded-xl border border-slate-200 bg-white p-5">
      <h3 className="font-semibold">{label}</h3>
      {src ? (
        // eslint-disable-next-line @next/next/no-img-element -- API uploads and local blob previews
        <img src={src} alt={file ? `Preview of new ${label.toLowerCase()}` : label} className="h-56 w-full rounded-lg bg-slate-100 object-contain" />
      ) : (
        <div className="flex h-40 items-center justify-center rounded-lg bg-slate-100 text-sm text-slate-500">No image set</div>
      )}
      <Input
        label={imagePath || savedPath ? "Replace image" : "Upload image"}
        type="file"
        accept=".jpg,.jpeg,.png,.webp,image/jpeg,image/png,image/webp"
        hint="JPG, PNG or WebP, up to 10 MB. Select an image to preview it, then save."
        disabled={saving}
        onChange={(e) => select(e.target.files?.[0] ?? null)}
      />
      {file && <p className="text-sm text-slate-600">Selected: {file.name}</p>}
      {error && <Alert tone="error">{error}</Alert>}
      {success && <Alert tone="success">Image saved. It is now visible on your public agency or tour page.</Alert>}
      <div className="flex gap-2">
        <Button type="submit" loading={saving} disabled={!file}>Save image</Button>
        {file && <Button type="button" variant="secondary" disabled={saving} onClick={clearSelection}>Discard selection</Button>}
      </div>
    </form>
  );
}
