export function roleNames(user = {}) {
  const roles = Array.isArray(user.roles) ? user.roles : (user.role ? [user.role] : [])
  return roles
    .map(role => typeof role === 'string' ? role : role?.name)
    .filter(Boolean)
}

export function isOrdinaryUser(user = {}) {
  const roles = roleNames(user)
  return roles.length === 0 || roles.every(role => role === 'user')
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
