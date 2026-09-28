import { imageUrl } from "@/lib/config";

/** Tour or agency image, with a neutral placeholder when there is none. */
export function TourImage({ path, alt, className = "" }: { path?: string | null; alt: string; className?: string }) {
  const src = imageUrl(path);
  if (!src) {
    return (
      <div className={`flex items-center justify-center bg-gradient-to-br from-teal-100 to-sky-100 text-3xl ${className}`} aria-hidden>
        🏞️
      </div>
    );
  }
  // Images are served by the API, so plain <img> avoids configuring the Next image optimizer for it.
  // eslint-disable-next-line @next/next/no-img-element
  return <img src={src} alt={alt} className={`object-cover ${className}`} />;
}
