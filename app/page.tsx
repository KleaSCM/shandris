"use client";

import { useState } from "react";
import InputBox from "@/components/InputBox";
import OutputBox from "@/components/OutputBox";
import styles from "@/styles/page.module.scss";

export default function Home() {
  const [response, setResponse] = useState("");

  return (
    <main className={styles.container}>
      <h1>Shandris</h1>
      <OutputBox output={response} />
      <InputBox setResponse={setResponse} />
      
    </main>
  );
}
