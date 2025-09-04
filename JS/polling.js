const apiEndpoint = 'https://api.github.com/repos/will-lp1/saru/issues';

const pollingInterval = 5000;
const maxPollingDuration = 300000;

async function pollApi(apiEndpoint, pollingInterval, maxPollingDuration) {
    const startTime = Date.now();

    const makeRequest = async () => {
        try {
            const response = await fetch(apiEndpoint);
            const data = await response.json();

           console.log('Fetched issues;', data.map(issue => issue.title));

            const elapsedTime = Date.now() - startTime;

            if (elapsedTime < maxPollingDuration){
                setTimeout(makeRequest, pollingInterval);
            } else {
                console.log('Maximum polling duraion reached. Stopping polling')
            }
        } catch (e) {
            console.error('Error making Api request:', e);
            const elapsedTime = Date.now() - startTime;

            if (elapsedTime <maxPollingDuration) {
                setTimeout(makeRequest, pollingInterval);
            } else {
                console.log("Max polling duaratin reached. stoppung polling")
            }
        }
    };
    makeRequest();
}

pollApi(apiEndpoint, pollingInterval, maxPollingDuration);
