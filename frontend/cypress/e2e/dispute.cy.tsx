import 'cypress-file-upload'
describe('Navigation', () => {
	//TODO: only login if not logged in
	let cookie : string | undefined
	it('should login', () => {
		cy.visit('https://capstone-dre.dns.net.za/login')
		cy.contains('Email').type('relanit981@digdy.com')
		cy.contains('Password').type('Password1234#')
		cy.get('button').contains('Login').click()
		cy.visit('http://capstone-dre.dns.net.za/disputes')
		cy.getCookie('jwt').then(tempCookie => {
			cookie = tempCookie?.value
		})
	})
	let time = Date.prototype.toUTCString()
	let title = 'Cypress Test Title' + time
	let description = 'I am creating a cypress test'
	it('should be able to create a dispute', () => {
		cy.setCookie('jwt', cookie ?? '')
		cy.visit('http://capstone-dre.dns.net.za/disputes')
		cy.get('a').contains('+ Create').click()
		cy.get('input[name="respondentName"]').click().type('Bob Charlie')
		cy.get('input[name="respondentEmail"]').click().type('bob@charlie.com')
		cy.get('input[name="respondentTelephone"]').click().type('0123456789')
		cy.get('input[name="title"]').click().type(title)
		cy.get('textarea[name="summary"').click().type(description)
		cy.fixture('testpicture.png').then(fileContent => {
			cy.get('input[type="file"]').attachFile({
				fileContent: fileContent.toString(),
				fileName: 'testpicture.png',
				mimeType: 'image/png'
			});
		});
		cy.get('button[type=submit]').click()
		cy.get('span').contains(title)
	})
	it('should be able see the dispute', () => {
		cy.setCookie('jwt', cookie ?? '')
		cy.visit('http://capstone-dre.dns.net.za/disputes')
		cy.get('span').contains(title).click()
		cy.get('h1').contains(title)
		cy.get('p').contains(description)
	})
  })