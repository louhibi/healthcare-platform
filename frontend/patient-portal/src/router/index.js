import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Patients from '../views/Patients.vue'
import DoctorAvailability from '../views/DoctorAvailability.vue'
import AppointmentManagement from '../views/AppointmentManagement.vue'
import AdminDashboard from '../views/AdminDashboard.vue'
import FormDemo from '../views/FormDemo.vue'
import DateLocalizationTest from '../components/DateLocalizationTest.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { 
      guest: true // Only accessible to non-authenticated users
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { 
      guest: true // Only accessible to non-authenticated users
    }
  },
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { 
      requiresAuth: true // Requires authentication
    }
  },
  {
    path: '/patients',
    name: 'Patients',
    component: Patients,
    meta: { 
      requiresAuth: true,
      roles: ['admin', 'doctor', 'nurse', 'staff'] // All authenticated users
    }
  },
  {
    path: '/availability',
    name: 'DoctorAvailability',
    component: DoctorAvailability,
    meta: { 
      requiresAuth: true,
      roles: ['admin', 'doctor'] // Only admins and doctors can manage availability
    }
  },
  {
    path: '/appointment-management',
    name: 'AppointmentManagement',
    component: AppointmentManagement,
    meta: { 
      requiresAuth: true,
      roles: ['admin', 'doctor', 'nurse', 'staff'] // All authenticated users
    }
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: AdminDashboard,
    meta: { 
      requiresAuth: true,
      roles: ['admin'] // Only admins can access admin dashboard
    }
  },
  {
    path: '/form-demo',
    name: 'FormDemo',
    component: FormDemo,
    meta: { 
      requiresAuth: true,
      roles: ['admin'] // Demo route for admins to test form configuration
    }
  },
  {
    path: '/date-test',
    name: 'DateLocalizationTest',
    component: DateLocalizationTest,
    meta: { 
      requiresAuth: true,
      roles: ['admin'] // Demo route for admins to test date localization
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Initialize auth store if not already done
  if (!authStore.isAuthenticated && localStorage.getItem('auth_token')) {
    try {
      await authStore.initialize()
    } catch (error) {
      console.error('Auth initialization failed:', error)
    }
  }

  const isAuthenticated = authStore.isAuthenticated
  const userRole = authStore.userRole

  // Check if route requires authentication
  if (to.meta.requiresAuth && !isAuthenticated) {
    // Store intended route for redirect after login
    localStorage.setItem('intended_route', to.fullPath)
    next('/login')
    return
  }

  // Check if route is for guests only (login/register)
  if (to.meta.guest && isAuthenticated) {
    // If user is already authenticated and going to login/register,
    // redirect to intended route or home
    const intendedRoute = localStorage.getItem('intended_route')
    if (intendedRoute) {
      localStorage.removeItem('intended_route')
      next(intendedRoute)
    } else {
      next('/')
    }
    return
  }

  // Check role-based access
  if (to.meta.roles && isAuthenticated) {
    if (!to.meta.roles.includes(userRole)) {
      // User doesn't have required role
      next('/')
      return
    }
  }

  next()
})

export default router