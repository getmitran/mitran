import { useState, useEffect, useRef } from 'react';

interface SSEEvent {
  type: string;
  data: string;
  timestamp: number;
}

export function useSSE(url = 'http://localhost:7780/api/v1/events') {
  const [events, setEvents] = useState<SSEEvent[]>([]);
  const [connected, setConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState<SSEEvent | null>(null);
  const esRef = useRef<EventSource | null>(null);

  useEffect(() => {
    const connect = () => {
      const es = new EventSource(url);
      esRef.current = es;

      es.onopen = () => setConnected(true);
      es.onerror = () => {
        setConnected(false);
        es.close();
        setTimeout(connect, 3000);
      };
      es.addEventListener('task_update', (e) => {
        const evt: SSEEvent = { type: 'task_update', data: e.data, timestamp: Date.now() };
        setLastEvent(evt);
        setEvents(prev => [...prev.slice(-99), evt]);
      });
      es.addEventListener('task_started', (e) => {
        const evt: SSEEvent = { type: 'task_started', data: e.data, timestamp: Date.now() };
        setLastEvent(evt);
        setEvents(prev => [...prev.slice(-99), evt]);
      });
      es.addEventListener('connected', () => setConnected(true));
    };

    connect();
    return () => { esRef.current?.close(); };
  }, [url]);

  return { events, connected, lastEvent };
}
