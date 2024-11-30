import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "VidioStream",
  description: "Best video streaming platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-br">
      <body>
        {children}
      </body>
    </html>
  );
}
