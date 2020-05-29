package main 

import (
    "io"
    "log"
    "os"
    "strconv"
    "strings"

    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"

    "gitlab.iqt.org/rashley/covid-test-db/models/poc"
    "gitlab.iqt.org/rashley/covid-test-db/models/company"
    "gitlab.iqt.org/rashley/covid-test-db/models/diagnostic"
    "gitlab.iqt.org/rashley/covid-test-db/models/diagnostic_type"
    "gitlab.iqt.org/rashley/covid-test-db/models/diagnostic_target_type"
    "gitlab.iqt.org/rashley/covid-test-db/models/regulatory_approval_type"
    "gitlab.iqt.org/rashley/covid-test-db/models/sample_type"
    "gitlab.iqt.org/rashley/covid-test-db/models/pcr_platform"
)

func molecularColumnMapping () map[string]int{
    return map[string]int{
        "company_name":  0,
        "test_name":   1,
        "pcr_platform": 2, 
        "sensitivity": 3,
        "specificity": 4,
        "source": 5,
        "test_url": 6,
        "regulatory_status": 7,
        "sample_type": 8,
        "point_of_care": 9,
        "prep_integrated": 10,
        "product_no": 11,
        "company_support": 12,
        "company_poc": 13,
        "company_poc_phone": 14,
        "company_street": 15,
        "company_city": 16,
        "company_state": 17,
        "company_country": 18,
        "company_postal_code": 19,
        "cost_per_kit": 20,
        "in_stock": 21,
        "lead_time": 22,
        "company_stage": 23,
        "company_valuation": 24,
        "company_entered_market": 25,
        "test_counts": 26,
        "production_rate": 27,
        "expanding_capacity": 28,
        "gene_target": 29,
        "verified_lod": 30,
        "avg_ct": 31,
        "turn_around_time": 32,
        "tests_per_run": 33,
        "tests_per_kit": 34,
    }
} 

func Index(vs []string, t string) int {
    for i, v := range vs {
        if v == t {
            return i
        }
    }
    return -1
}

func getDB () *gorm.DB {
    const addr = "postgresql://covid_bug@localhost:26257/covid_diagnostics?sslmode=disable"
    db, err := gorm.Open("postgres", addr)
    if err != nil {
        log.Fatal(err)
    }
    //defer db.Close()

    // Set to `true` and GORM will print out all DB queries.
    db.LogMode(false)

    return db
}

func getOrCreatePoc(name string, email string, phone string)(*poc.Poc, error){
    db := getDB()
    defer db.Close()
    var result *poc.Poc = nil
    existing, err := poc.FetchByNameAndEmail(db, name, email)
    if(existing != nil && !gorm.IsRecordNotFoundError(err)){
        result = existing
    } else {
        result, err = poc.Create(db, name, email, phone)
    }

    return result, err
}

func getOrCreateCompany(name string, streetAddress string, city string, state string, 
        country string, postalCode string, stage string, valuation string, poc poc.Poc)(*company.Company, error){
    db := getDB()
    defer db.Close()
    var result *company.Company = nil
    existing, err := company.FetchByName(db, name)
    if(existing != nil && !gorm.IsRecordNotFoundError(err)){
        result = existing
    } else {
        result, err = company.Create(db, name, streetAddress, city, state, country, postalCode, stage, valuation, poc)
    }

    return result, err
}

func getOrCreateTargetType(name string)(*diagnostic_target_type.DiagnosticTargetType, error){
    db := getDB()
    defer db.Close()
    var result *diagnostic_target_type.DiagnosticTargetType = nil
    existing, err := diagnostic_target_type.FetchByName(db, name)
    if(existing != nil && !gorm.IsRecordNotFoundError(err)){
        result = existing
    } else {
        result, err = diagnostic_target_type.Create(db, name)
    }

    return result, err
}

func getOrCreateSampleType(name string)(*sample_type.SampleType, error){
    db := getDB()
    defer db.Close()
    var result *sample_type.SampleType = nil
    existing, err := sample_type.FetchByName(db, name)
    if(existing != nil && !gorm.IsRecordNotFoundError(err)){
        result = existing
    } else {
        result, err = sample_type.Create(db, name)
    }

    return result, err
}

func getOrCreatePcrPlatform(name string)(*pcr_platform.PcrPlatform, error){
    db := getDB()
    defer db.Close()
    var result *pcr_platform.PcrPlatform = nil
    existing, err := pcr_platform.FetchByName(db, name)
    if(existing != nil && !gorm.IsRecordNotFoundError(err)){
        result = existing
    } else {
        result, err = pcr_platform.Create(db, name)
    }

    return result, err
}

func createDiagnostic( name string, description string, testUrl string, company company.Company, 
             diagnosticType diagnostic_type.DiagnosticType, poc poc.Poc, 
             verifiedLod string, avgCt float64, prepIntegrated bool,
             testsPerRun int64, testsPerKit int64, sensitivity float64, specificity float64,
             sourceOfPerfData string, catalogNo string, pointOfCare bool, costPerKit float64,
             inStock bool, leadTime int64,
             approvals []regulatory_approval_type.RegulatoryApprovalType, 
             targets []diagnostic_target_type.DiagnosticTargetType,
             sampleTypes []sample_type.SampleType,
             pcrPlatforms []pcr_platform.PcrPlatform)(*diagnostic.Diagnostic, error){
    db := getDB()
    defer db.Close()

    var result *diagnostic.Diagnostic = nil
    result, err := diagnostic.Create(db, name, description, testUrl, company, diagnosticType, poc,
                    verifiedLod, avgCt, prepIntegrated, testsPerRun, testsPerKit,
                    sensitivity, specificity, sourceOfPerfData, catalogNo, pointOfCare, costPerKit,
                    inStock, leadTime,
                    approvals, targets, sampleTypes, pcrPlatforms)   

    return result, err
}

func getSampleTypes(names []string)([]sample_type.SampleType, []error){
    db := getDB()
    defer db.Close()
    var types []sample_type.SampleType = nil
    var errs []error = nil

    for _, name := range names{
        st, err := getOrCreateSampleType(strings.TrimSpace(name))
        if(err == nil){
            types = append(types, *st)
        } else {
            errs = append(errs, err)
        }
    }

    return types, errs
}

func getPcrPlatforms(names []string)([]pcr_platform.PcrPlatform, []error){
    db := getDB()
    defer db.Close()
    var pcrs []pcr_platform.PcrPlatform = nil
    var errs []error = nil

    for _, name := range names{
        st, err := getOrCreatePcrPlatform(strings.TrimSpace(name))
        if(err == nil){
            pcrs = append(pcrs, *st)
        } else {
            errs = append(errs, err)
        }
    }

    return pcrs, errs
}

func getTargetTypes(names []string)([]diagnostic_target_type.DiagnosticTargetType, []error){
    db := getDB()
    defer db.Close()
    var types []diagnostic_target_type.DiagnosticTargetType
    var errs []error

    for _, name := range names{
        tt, err := getOrCreateTargetType(strings.TrimSpace(name))
        if(err == nil){
            types = append(types, *tt)
        } else {
            errs = append(errs, err)
        }
    }

    return types, errs
}

func getApprovals(name string)([]regulatory_approval_type.RegulatoryApprovalType, error){
    db := getDB()
    defer db.Close()
    var validApprovals []regulatory_approval_type.RegulatoryApprovalType
    allApprovals, err := regulatory_approval_type.FetchList(db)

    for _, a := range allApprovals{
        if(strings.Contains(strings.ToLower(name), strings.ToLower(a.Name))){
            validApprovals = append(validApprovals, a)
        }
    }

    return validApprovals, err
}

func getDiagnosticType(name string)(*diagnostic_type.DiagnosticType, error){
    db := getDB()
    defer db.Close()

    result, err := diagnostic_type.FetchByName(db, name)

    return result, err
}

func getDiagnosticFromRow(row []string)(*diagnostic.Diagnostic, error){
    var mapping = molecularColumnMapping()
    //PoC Data
    var contact_name string = strings.TrimSpace(row[mapping["company_poc"]])
    var contact_email string = strings.TrimSpace(row[mapping["company_support"]])
    var contact_phone string = strings.TrimSpace(row[mapping["company_poc_phone"]])

    poc, err := getOrCreatePoc(contact_name, contact_email, contact_phone)

    //company data
    company, err := getOrCreateCompany(
        strings.TrimSpace(row[mapping["company_name"]]),
        strings.TrimSpace(row[mapping["company_street"]]),
        strings.TrimSpace(row[mapping["company_city"]]),
        strings.TrimSpace(row[mapping["company_state"]]),
        strings.TrimSpace(row[mapping["company_country"]]),
        strings.TrimSpace(row[mapping["company_postal_code"]]),
        strings.TrimSpace(row[mapping["company_stage"]]),
        strings.TrimSpace(row[mapping["company_valuation"]]),
        *poc,
    )
    
    tts := strings.Split(row[mapping["gene_target"]], ",")
    sts := strings.Split(row[mapping["sample_type"]], ",")
    pcrs := strings.Split(row[mapping["pcr_platform"]], ",")
    targetTypes, errs := getTargetTypes(tts)
    if(len(errs) > 0){
        for _, e := range errs{
            log.Println(e)
        }
        return nil, errs[0]
    }
    sampleTypes, errs := getSampleTypes(sts)
    if(len(errs) > 0){
        for _, e := range errs{
            log.Println(e)
        }
        return nil, errs[0]
    }

    pcrPlatforms, errs := getPcrPlatforms(pcrs)
    if(len(errs) > 0){
        for _, e := range errs{
            log.Println(e)
        }
        return nil, errs[0]
    }
    
    diagnosticType, err := getDiagnosticType("molecular assays")
    approvals, err := getApprovals(strings.TrimSpace(row[mapping["regulatory_status"]]))

    if(err != nil){
        log.Println(err)
        return nil, err
    }
    positives := []string{"y", "yes", "t", "true", "1"}
    avgCt, _ := strconv.ParseFloat(row[mapping["avg_ct"]], 64)
    prepIntegrated := Index(positives, strings.ToLower(strings.TrimSpace(row[mapping["prep_integrated"]]))) >= 0
    tpr, _ := strconv.ParseInt(row[mapping["tests_per_run"]], 10, 64)
    tpk, _ := strconv.ParseInt(row[mapping["tests_per_kit"]], 10, 64)
    sensitivity, _ := strconv.ParseFloat(row[mapping["sensitivity"]], 64)
    specificity, _ := strconv.ParseFloat(row[mapping["specificity"]], 64)
    sourceOfPerfData := strings.TrimSpace(row[mapping["source"]])
    catalogNo := strings.TrimSpace(row[mapping["product_no"]])
    pointOfCare := Index(positives, strings.ToLower(strings.TrimSpace(row[mapping["point_of_care"]]))) >= 0
    costPerKit, _ := strconv.ParseFloat(row[mapping["Cost_per_kit"]], 64)
    inStock := Index(positives, strings.ToLower(strings.TrimSpace(row[mapping["in_stock"]]))) >= 0 
    leadTime, _ := strconv.ParseInt(row[mapping["lead_time"]], 10, 64)

    dx, dxErr := createDiagnostic(
        strings.TrimSpace(row[mapping["test_name"]]),
        strings.TrimSpace(row[mapping["test_name"]]),
        strings.TrimSpace(row[mapping["test_url"]]),
        *company,
        *diagnosticType,
        *poc,
        strings.TrimSpace(row[mapping["verified_lod"]]),
        avgCt,
        prepIntegrated,
        tpr,
        tpk,
        sensitivity,
        specificity,
        sourceOfPerfData,
        catalogNo,
        pointOfCare,
        costPerKit,
        inStock,
        leadTime,
        approvals,
        targetTypes,
        sampleTypes,
        pcrPlatforms,
    )

    return dx, dxErr
}

func main() {
    logFile, err := os.OpenFile("scraper.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        panic(err)
    }
    mw := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(mw)
    log.Println("Logging Started")
    log.Println("Excel file processing started")
    
    // open excel file
    f, err := excelize.OpenFile("Database_Molecular_and_Sero.xlsx")
    if err != nil {
        log.Println(err.Error())
        return
    }
    rows := f.GetRows("Molecular test fields")
    for _, row := range rows {
        _, dxErr := getDiagnosticFromRow(row)

        if(dxErr != nil){
            log.Println(dxErr)
        }
    }

    log.Println("Excel file processing complete")

}