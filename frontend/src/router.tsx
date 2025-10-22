import { useSyncExternalStore } from 'react'
import App from './App'
import AdminPanel from './pages/AdminPanel'

export function navigate(path: string) {
  if (window.location.pathname !== path) {
    window.history.pushState({}, '', path)
    window.dispatchEvent(new PopStateEvent('popstate'))
  }
}

function subscribe(callback: () => void) {
  window.addEventListener('popstate', callback)
  return () => window.removeEventListener('popstate', callback)
}

function getSnapshot() {
  return window.location.pathname
}

export function Router() {
  const pathname = useSyncExternalStore(subscribe, getSnapshot)

  if (pathname === '/adminpanel') {
    return <AdminPanel />
  }

  return <App />
}
