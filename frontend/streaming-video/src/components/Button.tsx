"use client";

import { Upload } from "lucide-react";
import { ReactNode } from "react";
import { useRouter } from "next/navigation";

export function SubmitButton({ isUploading }: { isUploading: boolean }) {
  return (
    <button
      type="submit"
      disabled={isUploading}
      className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-800 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span className="absolute left-0 inset-y-0 flex items-center pl-3">
        <Upload className="h-5 w-5 text-gray-500 group-hover:text-gray-400" aria-hidden="true" />
      </span>
      {isUploading ? "Enviando..." : "Enviar Vídeo"}
    </button>
  );
}

export function ButtonNavigate({ children, to, isUploading }: { children: ReactNode; to: string; isUploading?: boolean }) {

  const router = useRouter();
  
  return (
    <button disabled={isUploading} onClick={() => router.push(to)}>
      {children}
    </button>
  );
}