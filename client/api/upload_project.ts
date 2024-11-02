async function UploadProject(url: string, slug: string | null): Promise<any> {
    const requestBody = {
        gitURL: url,
        slug: slug
    };

    try {
        const response = await fetch('http://localhost:8080/project', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error uploading project:', error);
        throw error;
    }
}

export { UploadProject };