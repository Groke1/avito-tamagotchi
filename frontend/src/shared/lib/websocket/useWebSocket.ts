import { useGetWsTicketMutation } from '@/entities/pet'
import { SOCKET_URL, baseApi } from '@/shared/api/baseApi'
import { useAppDispatch, useAppSelector } from '@/shared/model/hooks'
import { useEffect, useRef } from 'react'
import { toast } from 'sonner'

export const useWebSocket = () => {
  const dispatch = useAppDispatch()
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)
  const [getWsTicket] = useGetWsTicketMutation()
  const socketRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    let isMounted = true

    if (!isAuthenticated) {
      if (socketRef.current) {
        socketRef.current.close()
        socketRef.current = null
      }
      return
    }

    const connectWS = async () => {
      try {
        const { ticket } = await getWsTicket().unwrap()
        if (!isMounted) return

        const ws = new WebSocket(`${SOCKET_URL}?ticket=${ticket}`)
        socketRef.current = ws

        ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            if (data.event_type === 'pet.updated' && data.payload) {
              dispatch(baseApi.util.invalidateTags(['Pet', 'Rewards']))
            } else if (data.event_type === 'leaderboard.position_updated') {
              dispatch(baseApi.util.invalidateTags(['Pet', 'User']))
            }
          } catch {
            toast.error('Ошибка при обработке сообщения сервера')
          }
        }

        ws.onclose = () => {
          socketRef.current = null

          if (isMounted && isAuthenticated) {
            reconnectTimeoutRef.current = setTimeout(connectWS, 3000)
          }
        }

        ws.onerror = () => {
          ws.close()
        }
      } catch {
        toast.error('Не удалось установить соединение с сервером')
        if (isMounted && isAuthenticated) {
          reconnectTimeoutRef.current = setTimeout(connectWS, 5000)
        }
      }
    }

    connectWS()

    return () => {
      isMounted = false
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }

      if (socketRef.current) {
        socketRef.current.close()
        socketRef.current = null
      }
    }
  }, [isAuthenticated, getWsTicket, dispatch])
}
