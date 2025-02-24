import React from "react";
import "../../public/globals.css"
import {Quicksand} from "next/font/google";

const quicksand = Quicksand({subsets: ['latin']})

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={quicksand.className}>
      <body>
        {children}
      </body>
    </html>
  );
}
