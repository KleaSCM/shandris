"use client";
import InputBox from '@/components/InputBox';
import OutputBox from '@/components/OutputBox';
import { useState } from 'react';

export default function Home() {
  const [response, setResponse] = useState('');

  return (
    <main>
      <h1>Shandris</h1>
      <InputBox />
      <OutputBox output={response} />
    </main>
  );
}
