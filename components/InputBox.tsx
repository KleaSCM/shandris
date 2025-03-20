"use client";

import { useState } from "react";
import ChatButton from "./ChatButton";
import styles from "@/styles/inputBox.module.scss";

export default function InputBox({ setResponse }: { setResponse: (text: string) => void }) {
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);

  const handleCopy = () => navigator.clipboard.writeText(input);
  const handlePaste = async () => setInput(await navigator.clipboard.readText());

  const handleSend = async () => {
    if (!input.trim()) return;

    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ prompt: input }),
      });

      const data = await res.json();
      setResponse(data.response);
    } catch (error) {
      console.error("Error sending request:", error);
      setResponse("Error: Could not get a response.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.inputContainer}>
      <textarea
        className={styles.inputBox}
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Type your message..."
      />
      <div className={styles.buttonGroup}>
        <ChatButton label="Copy" onClick={handleCopy} />
        <ChatButton label="Paste" onClick={handlePaste} />
        <ChatButton label={loading ? "Sending..." : "Send"} onClick={handleSend} disabled={loading} />
      </div>
    </div>
  );
}
