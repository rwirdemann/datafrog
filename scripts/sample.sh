for i in {1..20}
do
   echo "Running verfication $i..."
   curl --location --request PUT 'http://localhost:3000/tests/full.json/verifications'
   cd /Users/ralf/work/vscode/playwright-rt
   npx playwright test tests/full-10.spec.ts --project=chromium
   sleep 2
   curl --location --request DELETE 'http://localhost:3000/tests/full.json/verifications'
done

