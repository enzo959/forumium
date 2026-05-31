import { Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import NotFound from './pages/NotFound.tsx'
import ForgotPassword from './pages/ForgotPassword.tsx'
import ResetPassword from './pages/ResetPassword.tsx'
import Login from './pages/Login.tsx'
import Register from './pages/Register.tsx'
import PostDetail from './pages/PostDetail.tsx'
import PostForm from './pages/PostForm'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/posts/:id" element={<PostDetail />} />
      <Route path="forgot-password" element={<ForgotPassword />} />
      <Route path="reset-password" element={<ResetPassword />} />
      <Route path="login" element={<Login />} />
      <Route path="register" element={<Register />} />
      <Route path="/posts/:id" element={<PostDetail />} />
      <Route path="/posts/new" element={<PostForm />} />
      <Route path="/posts/:id/edit" element={<PostForm />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  )
}

export default App