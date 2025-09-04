
async function pollApi( ) {
    const startTime = Date.now();

    const makeRequest = async () => {
        try {
            const response = await fetch('https://api.github.com/repos/will-lp1/saru/issues');
            const data = await response.json();

           console.log('Fetched issues;', data.map(issue => issue.title));

            const elapsedTime = Date.now() - startTime;

            if (elapsedTime < 300000){
                setTimeout(makeRequest, 5000);
            } else {
                console.log('Maximum polling duraion reached. Stopping polling')
            }
        } catch (e) {
            console.error('Error making Api request:', e);
            const elapsedTime = Date.now() - startTime;

            if (elapsedTime < 300000) {
                setTimeout(makeRequest, 5000);
            } else {
                console.log("Max polling duaratin reached. stoppung polling")
            }
        }
    };
    makeRequest();
}

pollApi();
