import './App.css'
import { navigate } from './router'

function App() {

  return (
    <>
        <h2 className='heading'>Welcome to Admin Panel of Fake Salon!!</h2>
        <button className = 'panel-button' onClick={() => navigate('/adminpanel')}>Launch Admin panel</button>
    </>
  )
}

export default App
