
const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '/beacon', component: () => import('pages/Beacon.vue') },
      { path: '/', component: () => import('pages/Info.vue') },
      { path: '/devinfo', component: () => import('pages/Devinfo.vue') },
      { path: '/signup', component: () => import('pages/Signup.vue'), meta: { hideDrawer: true} }
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '*',
    component: () => import('pages/Error404.vue')
  }
]

export default routes
