import { useEffect, useState } from 'react'
import './AdminPanel.css'

type Row = {
  id: string | number
  query: string
  answer: string
  status: string | number
  customerId: string
  createdDate: Date
}

export default function AdminPanel() {
  const [rows, setRows] = useState<Row[]>([])
  const [isOpen, setIsOpen] = useState(false)
  const [selectedId, setSelectedId] = useState<string | number | null>(null)
  const [answerInput, setAnswerInput] = useState('')
  const [statusInput, setStatusInput] = useState<number>(0)
  const getStatus = (s: string | number) => {
    const n = typeof s === 'string' ? Number(s) : s
    switch (n) {
      case 1:
        return { text: 'Resolved', color: '#16a34a' } // green
      case 0:
        return { text: 'Pending', color: '#ca8a04' } // yellow
      case 2:
        return { text: 'Unresolved', color: '#dc2626' } // red
      default:
        return { text: String(s), color: '#6b7280' }
    }
  }

  const openModalFor = (row: Row) => {
    setSelectedId(row.id)
    setAnswerInput(row.answer ?? '')
    const n = typeof row.status === 'string' ? Number(row.status) : row.status
    setStatusInput([0,1,2].includes(Number(n)) ? Number(n) : 0)
    setIsOpen(true)
  }

  const closeModal = () => {
    setIsOpen(false)
    setSelectedId(null)
    setAnswerInput('')
    setStatusInput(0)
  }

  const saveChanges = async () => {
    if (selectedId == null) return
    try {
      const res = await fetch(`http://localhost:8080/api/v1/queries/${selectedId}/resolve`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ answer: answerInput, query_status: statusInput })
      })
      if (!res.ok) {
        return
      }
      setRows(prev => prev.map(r => r.id === selectedId ? { ...r, answer: answerInput, status: statusInput } : r))
      closeModal()
    } catch {
    }
  }

  useEffect(() => {
    (async () => {
      try {
        const res = await fetch('http://localhost:8080/api/v1/queries')
        if (!res.ok) {
          console.log(res);
          setRows([])
          return
        }
        const raw: any = await res.json()
        const list: any[] = Array.isArray(raw)
          ? raw
          : Array.isArray(raw?.data)
          ? raw.data
          : Array.isArray(raw?.results)
          ? raw.results
          : Array.isArray((raw as any)?.["Success: "])
          ? (raw as any)["Success: "]
          : []

        const mapped: Row[] = list.map((it: any) => ({
          id: it?.id ?? '',
          query: it?.query_text ?? '',
          answer: it?.answer ?? '',
          status: it?.query_status ?? '',
          customerId: it?.customer_id ?? '',
          createdDate: new Date(
            it?.created_at ?? Date.now()
          )
        }))

        setRows(mapped)
      } catch {
        setRows([])
      }
    })()
  }, [])

  return (
    <div>
      <h2 className="heading">Admin Panel</h2>
      <div style={{ overflowX: 'auto' }}>
        <table className="table">
          <thead>
            <tr>
              <th>Id</th>
              <th>Query</th>
              <th>Answer</th>
              <th>Status</th>
              <th>Customer Id</th>
              <th>Created Date</th>
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td colSpan={5} style={{ textAlign: 'center' }}>No data</td>
              </tr>
            ) : (
              rows.map((r) => {
                const st = getStatus(r.status)
                return (
                  <tr key={String(r.id)} onClick={() => openModalFor(r)} style={{ cursor: 'pointer' }}>
                    <td>{r.id}</td>
                    <td>{r.query}</td>
                    <td>{r.answer}</td>
                    <td>
                      <span style={{ color: st.color, fontWeight: 600 }}>{st.text}</span>
                    </td>
                    <td>{r.customerId}</td>
                    <td>{r.createdDate ? new Date(r.createdDate).toLocaleString() : ''}</td>
                  </tr>
                )
              })
            )}
          </tbody>
        </table>
      </div>
      {isOpen && (
        <div className="modal-backdrop">
          <div className="modal">
            <h3>Edit Answer</h3>
            <div>
              <label>
                Answer
                <textarea
                  className="modal-input"
                  value={answerInput}
                  onChange={(e) => setAnswerInput(e.target.value)}
                  rows={4}
                />
              </label>
            </div>
            <div>
              <label>
                Status
                <select
                  className="modal-select"
                  value={statusInput}
                  onChange={(e) => setStatusInput(Number(e.target.value))}
                >
                  <option value={1}>Resolved</option>
                  <option value={0}>Pending</option>
                  <option value={2}>Unresolved</option>
                </select>
              </label>
            </div>
            <div className="modal-actions">
              <button className="btn" onClick={closeModal}>Cancel</button>
              <button className="btn primary" onClick={saveChanges}>Save</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
