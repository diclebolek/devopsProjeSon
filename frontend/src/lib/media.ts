const DEFAULT_MEDIA_BASE = '/img';
const BLOB_HOST = 'insucomstorage.blob.core.windows.net';

export function getMediaBaseUrl(): string {
  return (process.env.NEXT_PUBLIC_MEDIA_BASE_URL || DEFAULT_MEDIA_BASE).replace(/\/$/, '');
}

function extractFilename(path: string): string {
  let filename = path.trim();

  if (filename.startsWith('http://') || filename.startsWith('https://')) {
    try {
      const parsed = new URL(filename);
      filename = parsed.pathname;
      const marker = '/medya/';
      const markerIndex = filename.indexOf(marker);
      if (markerIndex >= 0) {
        filename = filename.slice(markerIndex + marker.length);
      }
    } catch {
      return '';
    }
  }

  if (filename.startsWith('./')) filename = filename.slice(2);
  if (filename.startsWith('/assets/img/')) filename = filename.slice('/assets/img/'.length);
  else if (filename.startsWith('assets/img/')) filename = filename.slice('assets/img/'.length);
  else if (filename.startsWith('/img/')) filename = filename.slice('/img/'.length);
  else if (filename.startsWith('img/')) filename = filename.slice('img/'.length);
  filename = filename.replace(/^\/+/, '');

  return filename;
}

/**
 * Tema görselleri ve API medya alanları için URL üretir.
 * Azure Blob kapalıysa / yasaklıysa `public/img` altındaki dosyalara düşer.
 * Zaten yabancı bir http(s) adresi ise olduğu gibi döner.
 */
export function mediaUrl(path?: string | null): string {
  if (!path) return '';

  const trimmed = path.trim();
  if (!trimmed) return '';

  if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
    try {
      const parsed = new URL(trimmed);
      if (parsed.hostname !== BLOB_HOST) {
        return trimmed;
      }
    } catch {
      return trimmed;
    }
  }

  const filename = extractFilename(trimmed);
  if (!filename) return '';

  const encoded = filename.split('/').map((part) => encodeURIComponent(part)).join('/');
  return `${getMediaBaseUrl()}/${encoded}`;
}
