import "cypress-file-upload";
describe("Navigation", () => {
  let cookie: string | undefined;
  let time = Date.now();
  let title = "Cypress Test Title" + time;
  let description = "I am creating a cypress test";
  it("should be able to create a dispute", () => {
    cy.visit("/login");
    cy.viewport(1920, 2000)
    cy.contains("Email").type(Cypress.env('TEST_USER'));
    cy.contains("Password").type(Cypress.env('TEST_PASSWORD'));
    cy.get("button").contains("Login").click();
    cy.wait(200);
    cy.getCookie("jwt").then((tempCookie) => {
      cookie = tempCookie?.value;
    });
    cy.visit("/disputes");
    cy.visit("/disputes/create")
    cy.get('input[name="respondentName"]').click().type("Bob Charlie");
    cy.get('input[name="respondentEmail"]').click().type("yexiy79682@mvpalace.com");
    cy.get('input[name="respondentTelephone"]').click().type("0123456789");
    cy.get('input[name="title"]').click().type(title);
    cy.get('button').contains('Select a workflow').parent().click()
    cy.get('span').contains('New Workflow').click()
    cy.get('textarea[name="summary"').click().type(description);
    cy.fixture("test.txt").then((fileContent) => {
      cy.get('input[type="file"]').attachFile({
        fileContent: fileContent.toString(),
        fileName: "test.txt",
        mimeType: "image/png",
      });
    });
    cy.get('div').contains('Dispute Details').parent().parent().children().get('button').contains("Create").click()
  });
  it("should be able see the dispute", () => {
    cy.setCookie("jwt", cookie ?? "");
    cy.visit("/disputes");
    // cy.get("span").contains(title)//.click();
    // cy.get("h1").contains(title);
    // cy.get("p").contains(description);
  });
  it("should not be in the archive", () => {
    cy.setCookie("jwt", cookie ?? "");
    cy.visit("/archive");
    // cy.get('input[name="q"]').click().type(title);
    // cy.get('button[type="submit"]').click();
  });
});
