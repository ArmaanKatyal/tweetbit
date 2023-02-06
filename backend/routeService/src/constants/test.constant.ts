let testURI: string;

export const setURI = (uri: string) => {
    testURI = uri;
}

export const getURI = () => {
    console.log(testURI);
    return testURI;
}