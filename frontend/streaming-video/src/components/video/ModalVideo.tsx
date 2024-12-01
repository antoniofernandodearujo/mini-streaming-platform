import { useEffect, useRef, useState } from "react";
import { ModalTypes } from "@/types/modalTypes";
import Hls from "hls.js";
import { X } from "lucide-react";

export default function ModalVideo({ open, onClose, videoID }: ModalTypes) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const [videoURL, setVideoURL] = useState<string>("");

  useEffect(() => {
    if (!open) return; // Só execute se o modal estiver aberto

    const urlHls = `https://mini-streaming-platform-production.up.railway.app/videos/${videoID}`;

    const fetchVideoURL = async () => {
      try {
        const response = await fetch(urlHls);
        const data = await response.json();

        const defaultQuality = data.resolutions["480p"]?.[0]; // Primeiro item é o .m3u8 para 480p
        if (defaultQuality) {
          setVideoURL(defaultQuality); // Atualiza o estado com a URL da qualidade 480p
        } else {
          console.error("Qualidade de vídeo 480p não encontrada!");
        }
      } catch (error) {
        console.error("Erro ao obter a URL do vídeo:", error);
      }
    };

    fetchVideoURL();
  }, [open, videoID]);

  useEffect(() => {
    if (!videoURL || !videoRef.current) return;

    if (Hls.isSupported()) {
      const hls = new Hls({
        autoStartLoad: true, // Inicia a carga do vídeo automaticamente
        startLevel: 1, // Começar com 480p (nivel 1)
        capLevelToPlayerSize: true, // Limita a qualidade máxima à capacidade do player
        maxMaxBufferLength: 30, // Aumenta o tempo de buffering máximo para evitar mudanças de qualidade abruptas
      });

      hls.loadSource(videoURL); // Carregar o arquivo .m3u8
      hls.attachMedia(videoRef.current); // Anexar o vídeo ao HLS.js

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        videoRef.current?.play(); // Reproduzir o vídeo automaticamente
      });

      // Evento para verificar se a troca de qualidade está ocorrendo
      hls.on(Hls.Events.LEVEL_SWITCHED, (event, data) => {
        console.log("Nível de qualidade trocado para:", data.level);
      });

      // Limpeza ao desmontar o componente ou alterar o videoURL
      return () => {
        hls.destroy();
      };
    } else if (videoRef.current.canPlayType("application/vnd.apple.mpegurl")) {
      // Caso HLS não seja suportado, mas o navegador tenha suporte nativo
      videoRef.current.src = videoURL;
      videoRef.current.addEventListener("loadedmetadata", () => {
        videoRef.current?.play();
      });
    }
  }, [videoURL]);

  return (
    <div
      className={`fixed inset-0 bg-black bg-opacity-70 z-50 ${
        open ? "block" : "hidden"
      }`}
    >
      <div className="fixed inset-0 flex justify-center items-center">
        <div className="bg-gray-900 rounded-lg p-4 w-96 h-96 relative">
          <video
            ref={videoRef}
            controls
            className="w-full h-full"
          ></video>
          <button
            onClick={onClose}
            className="absolute top-2 right-2 bg-red-600 text-white p-2 rounded-full hover:bg-red-700"
          >
            <X size={22} />
          </button>
        </div>
      </div>
    </div>
  );
}
