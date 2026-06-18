import { useEffect, useRef, useState, useCallback } from 'react'

type EventType = 'task.created' | 'task.updated' | 'task.checkpoint' | 'agent.status' | 'chat.message'

interface WSEvent {
  type: EventType
  data: any
}

type EventHandler = (event: WSEvent) => void

export function useWebSocket(url: string = 'ws://localhost:7780/api/v1/ws') {
  const [connected, setConnected] = useState(false)
  const [lastEvent, setLastEvent] = useState<WSEvent | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const handlersRef = useRef<Map<EventType, EventHandler[]>>(new Map())
  const reconnectRef = useRef<number>(0)

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    const ws = new WebSocket(url)

    ws.onopen = () => {
      setConnected(true)
      reconnectRef.current = 0
    }

    ws.onmessage = (e) => {
      try {
        const event: WSEvent = JSON.parse(e.data)
        setLastEvent(event)
        const handlers = handlersRef.current.get(event.type) || []
        handlers.forEach(h => h(event))
      } catch {}
    }

    ws.onclose = () => {
      setConnected(false)
      const delay = Math.min(1000 * 2 ** reconnectRef.current, 30000)
      reconnectRef.current++
      setTimeout(connect, delay)
    }

    ws.onerror = () => ws.close()
    wsRef.current = ws
  }, [url])

  useEffect(() => {
    connect()
    return () => { wsRef.current?.close() }
  }, [connect])

  const on = useCallback((type: EventType, handler: EventHandler) => {
    const handlers = handlersRef.current.get(type) || []
    handlers.push(handler)
    handlersRef.current.set(type, handlers)
    return () => {
      const h = handlersRef.current.get(type) || []
      handlersRef.current.set(type, h.filter(x => x !== handler))
    }
  }, [])

  const send = useCallback((data: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data))
    }
  }, [])

  return { connected, lastEvent, on, send }
}
