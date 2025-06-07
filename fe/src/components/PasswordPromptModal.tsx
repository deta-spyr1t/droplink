import { useState } from "react";

interface PasswordPromptModalProps {
  onSubmit: (password: string) => void;
  triesLeft: number;
  disabled: boolean;
}

export function PasswordPromptModal({ onSubmit, triesLeft, disabled }: PasswordPromptModalProps) {
  const [password, setPassword] = useState("");

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Enter Password to Decrypt File</h2>
        <input
          type="password"
          disabled={disabled}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Enter password"
          autoFocus
        />
        <button
          disabled={disabled || !password}
          onClick={() => {
            onSubmit(password);
            setPassword("");
          }}
        >
          Decrypt & Download
        </button>
        <p>{disabled ? "Please wait..." : `Tries left: ${triesLeft}`}</p>
      </div>
    </div>
  );
}
