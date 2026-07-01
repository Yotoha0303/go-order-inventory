import { AxiosError } from 'axios'
import { toast } from 'sonner'

export function handleServerError(error: unknown) {
  if (import.meta.env.DEV) {
    // eslint-disable-next-line no-console
    console.log(error)
  }

  let errMsg = 'Something went wrong!'

  if (
    error &&
    typeof error === 'object' &&
    'status' in error &&
    Number(error.status) === 204
  ) {
    errMsg = 'No content.'
  }

  if (error instanceof AxiosError) {
    const data = error.response?.data

    if (data && typeof data === 'object') {
      if ('message' in data && typeof data.message === 'string') {
        errMsg = data.message
      } else if ('title' in data && typeof data.title === 'string') {
        errMsg = data.title
      }
    }
  } else if (error instanceof Error && error.message) {
    errMsg = error.message
  }

  toast.error(errMsg)
}
