import { AuthLayout } from '@/app/layouts/AuthLayout'
import { RegisterForm } from '@/features/auth'

const Register = () => {
  return (
    <AuthLayout>
      <RegisterForm />
    </AuthLayout>
  )
}

export default Register
