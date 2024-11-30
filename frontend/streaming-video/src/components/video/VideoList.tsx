"use client";

import { useEffect, useState } from "react";
import { PlayIcon } from "lucide-react";
import ModalVideo from "./ModalVideo";
import { getThumbnail } from "@/functions/videoList";
import ListAllVideos from "@/core/ListAllVideos";
import { VideoTypes } from "@/types/videoTypes";
import { ButtonNavigate } from "@/components/Button";
import Loading from "../Loading";

export default function VideoList() {
  const [videos, setVideos] = useState<VideoTypes[]>([]); // Inicializa como um array vazio
  const [selectedVideo, setSelectedVideo] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(true);
  
  // Carregar vídeos ao montar o componente
  useEffect(() => {

    setIsLoading(true);
    ListAllVideos.listAllVideos().then((response) => {
      if (response.success) {
        setVideos(response.videos);
      }
    });

    setIsLoading(false);

    // Função de limpeza
    return () => {
      setVideos([]);
    };
  }, []);

  const handleSelectVideo = (videoID: string) => {
    console.log("Vídeo clicado:", videoID);  // Verifique o videoID aqui
    if (!videoID) {
      console.error("ID do vídeo não fornecido.");
      return;
    }
    setSelectedVideo(videoID);
    console.log("Estado de selectedVideo após set:", { videoID });
  };


  return (
    <>
      {isLoading ? (
        <Loading isLoading={isLoading} />
      ) : (
        <div className="p-4 flex justify-center items-center flex-col">
          <h1 className="text-2xl font-bold text-white mb-4">Lista de Vídeos</h1>
          <div className="mt-3 mb-3 text-center text-lg rounded-md px-5 py-2 bg-gray-800 font-bold hover:opacity-80 transition-all duration-150 text-white">
            <ButtonNavigate to="/">Enviar meu Vídeo</ButtonNavigate>
          </div>
          <hr className="mb-4 h-1 bg-white w-11/12 border-none" />
          <ul className="flex flex-row flex-wrap gap-8 justify-center items-center">
            {videos.map((video, index) => (
              <li key={index} className="relative bg-gray-900 rounded-lg overflow-hidden transform transition-transform w-80 h-72 hover:scale-105">
                <img
                  src={getThumbnail(video)} // Usando a função getThumbnail para a URL
                  alt={`Thumbnail do vídeo ${video}`} // Alt para acessibilidade
                  className="w-full h-full object-cover" // Ajusta a imagem para cobrir toda a área
                />
                <div className="absolute inset-0 flex justify-center items-center z-10 text-white">
                  <button
                    className="bg-red-600 text-white p-4 rounded-full hover:bg-red-700 transition-colors"
                    onClick={() => {
                      console.log("Vídeo clicado:", video);  // Inspecione o objeto `video`
                      if (video) {
                        handleSelectVideo(video); // Passa o videoID corretamente
                      } else {
                        console.error("videoID não encontrado para o vídeo:", video);
                      }
                    }}
                  >
                    <PlayIcon size={32} />
                  </button>
                </div>
              </li>
            ))}
          </ul>

          {/* Modal de vídeo */}
        {selectedVideo && (
          <ModalVideo
            open={!!selectedVideo}
            onClose={() => setSelectedVideo('')} // Fecha o modal ao clicar no botão de fechar
            videoID={selectedVideo} // Passa o videoID diretamente
            qualities={["1080p", "720p", "480p"]}
          />
        )}
      </div>
      )}
    </>
  );
}
