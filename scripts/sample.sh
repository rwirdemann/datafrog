echo "Running $1 verifications."
for i in $(seq 1 $1); do
   echo "Verfication $i..."
   curl --location --request PUT 'http://localhost:3000/tests/full.json/verifications'
   cd /Users/ralf/work/vscode/playwright-rt
   npx playwright test tests/full-10.spec.ts --project=chromium
   sleep 1
   curl --location --request DELETE 'http://localhost:3000/tests/full.json/verifications'
done

