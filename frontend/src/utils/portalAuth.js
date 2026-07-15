export function roleNames(user = {}) {
  const roles = Array.isArray(user.roles) ? user.roles : (user.role ? [user.role] : [])
  return roles
    .map(role => typeof role === 'string' ? role : role?.name)
    .filter(Boolean)
}

export function isOrdinaryUser(user = {}) {
  const roles = roleNames(user)
  // Administrative/managerial roles go to /portal; athletes, coaches, officials, club managers, and standard users go to /dashboard.
  const adminRoles = [
    'admin', 'super_admin', 'general_secretary', 'federation_officer', 'portal_operator',
    'role_admin', 'role_super_admin', 'role_general_secretary', 'role_federation_officer', 'role_portal_operator'
  ]
  return roles.length === 0 || !roles.some(role => adminRoles.includes(role))
}

export function portalDestination(user = {}) {
  return isOrdinaryUser(user) ? '/dashboard' : '/portal'
}

export function storedPortalUser() {
  try {
    return JSON.parse(localStorage.getItem('ncsms_user') || '{}')
  } catch {
    return {}
  }
}
