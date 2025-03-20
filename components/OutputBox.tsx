"use client";

import ChatButton from "./ChatButton";
import styles from "@/styles/outputBox.module.scss";

interface OutputBoxProps {
  output: string;
}

export default function OutputBox({ output }: OutputBoxProps) {
  const handleCopy = () => navigator.clipboard.writeText(output);

  return (
    <div className={styles.outputContainer}>
      <div className={styles.outputBox}>{output || "Shandris..."}</div>
      <div className={styles.buttonGroup}>
        <ChatButton label="Copy" onClick={handleCopy} />
      </div>
    </div>
  );
}
