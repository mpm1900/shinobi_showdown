import { setResponseHeader } from '@tanstack/react-start/server'

function setResponseCookie(response: Response) {
  const setCookies =
    (
      response.headers as Headers & { getSetCookie?: () => string[] }
    ).getSetCookie?.() ?? []
  const setCookie = setCookies[0] ?? response.headers.get('set-cookie')
  if (setCookie && setCookie.length > 0) {
    setResponseHeader('set-cookie', setCookie)
  }
}

export { setResponseCookie }
