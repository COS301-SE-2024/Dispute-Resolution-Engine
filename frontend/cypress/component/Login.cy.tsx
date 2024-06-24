describe('Archive Dispute Details', () => {
  it('successfully retrieves an archived dispute by ID', () => {
    cy.request('GET', 'http://localhost:8080/disputes/archive/1').then((response) => {
      console.log(response)
      expect(response.body.data).to.have.property('id', 1)
    });
  });
});