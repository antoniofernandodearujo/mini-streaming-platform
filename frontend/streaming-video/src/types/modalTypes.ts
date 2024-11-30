export interface ModalTypes {
    open: boolean;
    onClose: () => void;
    videoID: string;
    qualities: ('1080p' | '720p' | '480p')[];
}