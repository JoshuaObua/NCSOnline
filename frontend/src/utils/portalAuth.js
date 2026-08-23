export function roleNames(user = {}) {
  const roles = Array.isArray(user.roles) ? user.roles : (user.role ? [user.role] : [])
  return roles
    .map(role => typeof role === 'string' ? role : role?.name)
    .filter(Boolean)
}

export function isOrdinaryUser(user = {}) {
  const roles = roleNames(user)
  // Administrative/managerial and sports registry officer roles go to /portal; pure applicants/standard users go to /dashboard.
  const adminRoles = [
    'admin', 'super_admin', 'general_secretary', 'federation_officer', 'portal_operator',
    'medical_officer', 'anti_doping_officer', 'technical_official', 'club_manager', 'coach',
    'role_admin', 'role_super_admin', 'role_general_secretary', 'role_federation_officer', 'role_portal_operator',
    'role_medical_officer', 'role_anti_doping_officer', 'role_technical_official'
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
