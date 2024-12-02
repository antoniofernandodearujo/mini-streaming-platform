import axios from "axios";

class VideoLister {
    private baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    // Lista todos os vídeos disponíveis
    async listAllVideos(): Promise<{ success: boolean; videos: any[] }> {
        try {
            const response = await axios.get(`${this.baseURL}/videos`);

            if (response.status === 200) {
                return { success: true, videos: response.data };
            } else {
                return { success: false, videos: [] };
            }
        } catch (error) {
            console.error("Erro ao listar vídeos:", error);
            return { success: false, videos: [] };
        }
    }
}
// Exporta uma instância já configurada da classe
export default new VideoLister("http://localhost:8080");
