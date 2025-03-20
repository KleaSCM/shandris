import type { Metadata } from 'next';
import './globals.scss';


export const metadata: Metadata = {
  title: 'Shandris',
  description: 'this does things',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
