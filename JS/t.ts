import axios from "axios";

const sleep = (ms: number) =>
    new Promise(resolve => setTimeout(resolve, ms));

function isRetryableError(error: unknown): boolean {
    if (!axios.isAxiosError(error)) {
        return false;
    }

    // Timeout
    if (
        error.code === "ECONNABORTED" ||
        error.code === "ETIMEDOUT"
    ) {
        return true;
    }

    // Network-level failure
    if (error.code === "ERR_NETWORK") {
        return true;
    }

    // Transient HTTP failures
    const status = error.response?.status;

    return status === 502 ||
           status === 503 ||
           status === 504;
}

export async function retry<T>(
    operation: () => Promise<T>,
    maxAttempts = 3
): Promise<T> {

    for (let attempt = 1; attempt <= maxAttempts; attempt++) {

        try {
            return await operation();

        } catch (error) {

            const retryable = isRetryableError(error);

            // Don't retry permanent/non-transient errors
            if (!retryable) {
                throw error;
            }

            // We've used all attempts
            if (attempt === maxAttempts) {
                throw error;
            }

            const delay = 100 * Math.pow(2, attempt - 1);

            console.log(
                `Request failed. Retrying in ${delay}ms... ` +
                (attempt ${attempt + 1}/${maxAttempts})
            );

            await sleep(delay);
        }
    }

    throw new Error("Retry failed unexpectedly");