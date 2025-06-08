// components/DownloadLinkWithCopy.tsx
import { useState } from "react";

interface DownloadLinkWithCopyProps {
  downloadUrl: string;
}

export function CopyClick({ downloadUrl }: DownloadLinkWithCopyProps) {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(downloadUrl).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <div className="flex items-center space-x-2 mt-2">
      <button
        onClick={copyToClipboard}
        className="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300 text-sm"
        aria-label="Copy download link"
      >
        {copied ? "Link copied!" : "Copy Link"}
      </button>
    </div>
  );
}
