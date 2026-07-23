/**
 * Enterprise Report Exporter Utility
 * Supports exporting structured datasets to CSV, Excel (.xls/.xlsx), Word (.doc/.docx), and PDF.
 */

export function exportToCSV(filename, headers, rows) {
  const processCell = (val) => {
    if (val === null || val === undefined) return '""'
    const str = String(val).replace(/"/g, '""')
    return `"${str}"`
  }

  const csvContent = [
    headers.map(h => processCell(h.label)).join(','),
    ...rows.map(row => headers.map(h => processCell(row[h.key])).join(','))
  ].join('\r\n')

  // UTF-8 BOM for Excel compatibility
  const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' })
  triggerDownload(blob, `${filename}.csv`)
}

export function exportToExcel(filename, headers, rows, title = 'Official Report') {
  const tableHead = headers.map(h => `<th style="background-color: #1e3a8a; color: #ffffff; font-weight: bold; padding: 8px; border: 1px solid #cbd5e1;">${h.label}</th>`).join('')
  const tableRows = rows.map(r => {
    const cells = headers.map(h => `<td style="padding: 6px; border: 1px solid #cbd5e1;">${r[h.key] !== undefined && r[h.key] !== null ? r[h.key] : ''}</td>`).join('')
    return `<tr>${cells}</tr>`
  }).join('')

  const excelTemplate = `
    <html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
    <head>
      <meta charset="utf-8" />
      <!--[if gte mso 9]>
      <xml>
        <x:ExcelWorkbook>
          <x:ExcelWorksheets>
            <x:ExcelWorksheet>
              <x:Name>${title}</x:Name>
              <x:WorksheetOptions>
                <x:DisplayGridlines/>
              </x:WorksheetOptions>
            </x:ExcelWorksheet>
          </x:ExcelWorksheets>
        </x:ExcelWorkbook>
      </xml>
      <![endif]-->
      <style>
        body { font-family: Arial, sans-serif; font-size: 11pt; }
        table { border-collapse: collapse; width: 100%; }
        th { background-color: #1e3a8a; color: #ffffff; text-align: left; }
        tr:nth-child(even) { background-color: #f8fafc; }
      </style>
    </head>
    <body>
      <h2>NATIONAL COUNCIL OF SPORTS - ${title.toUpperCase()}</h2>
      <p>Report Generated: ${new Date().toLocaleString()} | Total Records: ${rows.length}</p>
      <table>
        <thead><tr>${tableHead}</tr></thead>
        <tbody>${tableRows}</tbody>
      </table>
    </body>
    </html>
  `

  const blob = new Blob([excelTemplate], { type: 'application/vnd.ms-excel;charset=utf-8' })
  triggerDownload(blob, `${filename}.xls`)
}

export function exportToDocx(filename, headers, rows, title = 'Official Report') {
  const tableHead = headers.map(h => `<th style="background-color: #1e3a8a; color: #ffffff; font-weight: bold; padding: 8px; border: 1px solid #94a3b8;">${h.label}</th>`).join('')
  const tableRows = rows.map((r, idx) => {
    const bg = idx % 2 === 0 ? '#ffffff' : '#f8fafc'
    const cells = headers.map(h => `<td style="padding: 6px; border: 1px solid #cbd5e1; background-color: ${bg};">${r[h.key] !== undefined && r[h.key] !== null ? r[h.key] : ''}</td>`).join('')
    return `<tr>${cells}</tr>`
  }).join('')

  const docContent = `
    <html xmlns:o='urn:schemas-microsoft-com:office:office' xmlns:w='urn:schemas-microsoft-com:office:word' xmlns='http://www.w3.org/TR/REC-html40'>
    <head>
      <meta charset='utf-8'>
      <title>${title}</title>
      <style>
        @page { size: A4 landscape; margin: 1in; }
        body { font-family: 'Calibri', Arial, sans-serif; font-size: 10pt; color: #0f172a; }
        h1 { color: #1e3a8a; font-size: 16pt; margin-bottom: 2px; }
        h2 { color: #334155; font-size: 12pt; margin-top: 0; margin-bottom: 15px; text-transform: uppercase; }
        table { border-collapse: collapse; width: 100%; margin-top: 15px; }
        .meta { color: #64748b; font-size: 9pt; margin-bottom: 20px; }
      </style>
    </head>
    <body>
      <h1>NATIONAL COUNCIL OF SPORTS</h1>
      <h2>${title}</h2>
      <div class="meta">
        <strong>Generated Date:</strong> ${new Date().toLocaleString()} &nbsp;|&nbsp; 
        <strong>Total Items:</strong> ${rows.length}
      </div>
      <table border="1">
        <thead><tr>${tableHead}</tr></thead>
        <tbody>${tableRows}</tbody>
      </table>
    </body>
    </html>
  `

  const blob = new Blob([docContent], { type: 'application/msword;charset=utf-8' })
  triggerDownload(blob, `${filename}.doc`)
}

export function exportToPDF(title, headers, rows, metadata = {}) {
  const tableHead = headers.map(h => `<th style="background-color: #1e3a8a; color: #ffffff; font-weight: bold; padding: 6px 8px; border: 1px solid #1e3a8a; font-size: 9pt;">${h.label}</th>`).join('')
  const tableRows = rows.map((r, idx) => {
    const bg = idx % 2 === 0 ? '#ffffff' : '#f8fafc'
    const cells = headers.map(h => `<td style="padding: 6px 8px; border: 1px solid #cbd5e1; background-color: ${bg}; font-size: 8.5pt;">${r[h.key] !== undefined && r[h.key] !== null ? r[h.key] : ''}</td>`).join('')
    return `<tr>${cells}</tr>`
  }).join('')

  const printWindow = window.open('', '_blank')
  if (!printWindow) return alert('Please allow popups to generate PDF report.')

  const html = `
    <!DOCTYPE html>
    <html>
    <head>
      <title>${title} - NCS Uganda</title>
      <style>
        @page { size: A4 landscape; margin: 12mm; }
        body { font-family: Arial, Helvetica, sans-serif; color: #0f172a; margin: 0; padding: 0; }
        .header { display: flex; align-items: center; justify-content: space-between; border-bottom: 3px solid #1e3a8a; padding-bottom: 10px; margin-bottom: 15px; }
        .org-title { font-size: 16pt; font-weight: 800; color: #1e3a8a; margin: 0; }
        .report-title { font-size: 13pt; font-weight: 700; color: #334155; margin-top: 4px; margin-bottom: 0; text-transform: uppercase; }
        .meta-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; background: #f1f5f9; padding: 8px 12px; border-radius: 4px; margin-bottom: 15px; font-size: 8.5pt; }
        .meta-item strong { color: #475569; display: block; text-transform: uppercase; font-size: 7.5pt; }
        table { border-collapse: collapse; width: 100%; page-break-inside: auto; }
        tr { page-break-inside: avoid; page-break-after: auto; }
        th { text-transform: uppercase; text-align: left; }
        footer { margin-top: 20px; font-size: 8pt; color: #94a3b8; text-align: center; border-top: 1px solid #e2e8f0; padding-top: 8px; }
      </style>
    </head>
    <body>
      <div class="header">
        <div>
          <h1 class="org-title">NATIONAL COUNCIL OF SPORTS</h1>
          <h2 class="report-title">${title}</h2>
        </div>
        <div style="text-align: right; font-size: 8pt; color: #64748b;">
          <strong>Date:</strong> ${new Date().toLocaleDateString()}<br/>
          <strong>Time:</strong> ${new Date().toLocaleTimeString()}
        </div>
      </div>

      <div class="meta-grid">
        <div class="meta-item"><strong>Date Range</strong> ${metadata.dateRange || 'All Time'}</div>
        <div class="meta-item"><strong>Department Filter</strong> ${metadata.department || 'All Departments'}</div>
        <div class="meta-item"><strong>Status Filter</strong> ${metadata.status || 'All Statuses'}</div>
        <div class="meta-item"><strong>Total Records</strong> ${rows.length} Entries</div>
      </div>

      <table>
        <thead><tr>${tableHead}</tr></thead>
        <tbody>${tableRows}</tbody>
      </table>

      <footer>
        Official Administrative System Report &bull; National Council of Sports, Lugogo Sports Complex, Kampala &bull; Confidential
      </footer>

      <script>
        window.onload = function() {
          window.print();
          setTimeout(function() { window.close(); }, 750);
        }
      </script>
    </body>
    </html>
  `

  printWindow.document.open()
  printWindow.document.write(html)
  printWindow.document.close()
}

function triggerDownload(blob, filename) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
