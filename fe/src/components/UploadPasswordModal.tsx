import { useState } from "react";

interface UploadPasswordModalProps {
  onConfirm: (password: string) => void;
  onCancel: () => void;
}

export function UploadPasswordModal({ onConfirm, onCancel }: UploadPasswordModalProps) {
  const [password, setPassword] = useState("");

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Enter Password to Encrypt File</h2>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Enter a strong password"
          autoFocus
        />
        <div className="modal-buttons">
          <button onClick={() => onCancel()}>Cancel</button>
          <button disabled={!password} onClick={() => onConfirm(password)}>
            Upload
          </button>
        </div>
      </div>
    </div>
  );
}
