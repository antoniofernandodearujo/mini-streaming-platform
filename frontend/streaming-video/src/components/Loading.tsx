interface LoadingProps {
    isLoading: boolean;
}

export default function(isLoading: LoadingProps) {
    if (!isLoading) return null;

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <div className="animate-spin rounded-full h-32 w-32 border-t-4 border-b-4 border-blue-500"></div>
        </div>
    );
};
