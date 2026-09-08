export const expenseRoles = ['super_admin', 'admin', 'accountant', 'senior_accountant', 'chief_accountant', 'asset_accountant', 'finance_department']
export const expenseData = res => res?.data?.data ?? res?.data
export const expenseError = error => error.response?.data?.error?.message || error.response?.data?.message || 'Unable to complete the request. Please try again.'
export const expenseMoney = amount => new Intl.NumberFormat('en-UG', {minimumFractionDigits:2,maximumFractionDigits:2}).format(Number(amount || 0))
export const paymentMethods = {CASH:'Cash',BANK_TRANSFER:'Bank transfer',MOBILE_MONEY:'Mobile money',CARD:'Card',OTHER:'Other'}
