'use client';

import { useState } from 'react';
import ChatButton from './ChatButton';
import styles from '@/styles/inputBox.module.scss';

export default function InputBox() {
  const [input, setInput] = useState('');

  const handleCopy = () => {
    navigator.clipboard.writeText(input);
  };

  const handlePaste = async () => {
    const text = await navigator.clipboard.readText();
    setInput(text);
  };

  return (
    <div className={styles.inputContainer}>
      <textarea
        className={styles.inputBox}
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Type your message here..."
      />
      <div className={styles.buttonGroup}>
        <ChatButton label="Copy" onClick={handleCopy} />
        <ChatButton label="Paste" onClick={handlePaste} />
      </div>
    </div>
  );
}
