import { useState, useEffect } from 'react'
import { useWebSocket } from './useWebSocket'
import { usePolling } from './usePolling'
import { api } from '../api'

export function useRealtimeTasks() {
  const { connected, on } = useWebSocket()
  const polling = usePolling(api.fetchTasks, connected ? 30000 : 3000)
  const [tasks, setTasks] = useState<any[]>([])

  useEffect(() => {
    if (polling.data) setTasks(polling.data)
  }, [polling.data])

  useEffect(() => {
    if (!connected) return
    const unsub1 = on('task.created', (e) => setTasks(prev => [...prev, e.data]))
    const unsub2 = on('task.updated', (e) => setTasks(prev => prev.map(t => t.id === e.data.id ? e.data : t)))
    return () => { unsub1(); unsub2() }
  }, [connected, on])

  return { tasks, connected, refetch: polling.refetch, loading: polling.loading }
}

export function useRealtimeAgents() {
  const { connected, on } = useWebSocket()
  const polling = usePolling(api.fetchAgents, connected ? 60000 : 5000)
  const [agents, setAgents] = useState<any[]>([])

  useEffect(() => {
    if (polling.data) setAgents(polling.data)
  }, [polling.data])

  useEffect(() => {
    if (!connected) return
    return on('agent.status', (e) => setAgents(prev => prev.map(a => a.id === e.data.id ? e.data : a)))
  }, [connected, on])

  return { agents, connected, loading: polling.loading }
}
