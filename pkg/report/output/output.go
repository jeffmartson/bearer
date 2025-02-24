package output

import (
	"errors"
	"fmt"

	"github.com/bearer/bearer/pkg/commands/process/filelist/files"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/basebranchfindings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/saas"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/types"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func GetOutput(
	report types.Report,
	config settings.Config,
	files []files.File,
	baseBranchFindings *basebranchfindings.Findings,
) (data any, flow *dataflow.DataFlow, reportPassed bool, err error) {
	reportPassed = true

	switch config.Report.Report {
	case flag.ReportDetectors:
		data, flow, err = detectors.GetOutput(report, config)
	case flag.ReportDataFlow:
		data, flow, err = GetDataflow(report, config, false)
	case flag.ReportSecurity:
		data, flow, reportPassed, err = reportSecurity(report, config, files, baseBranchFindings)
	case flag.ReportSaaS:
		securityResults, dataflow, _, secErr := reportSecurity(report, config, files, baseBranchFindings)

		if secErr != nil {
			err = secErr
			return
		}

		data, flow, err = saas.GetReport(config, securityResults, dataflow, files)
	case flag.ReportPrivacy:
		data, flow, err = getPrivacyReportOutput(report, config)
	case flag.ReportStats:
		data, flow, err = reportStats(report, config)
	default:
		err = fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
	}
	return
}

func GetPrivacyReportCSVOutput(report types.Report, dataflow *dataflow.DataFlow, config settings.Config) (*string, error) {
	csvString, err := privacy.BuildCsvString(dataflow, config)
	if err != nil {
		return nil, err
	}

	content := csvString.String()

	return &content, nil
}

func getPrivacyReportOutput(report types.Report, config settings.Config) (*privacy.Report, *dataflow.DataFlow, error) {
	dataflow, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, err
	}

	return privacy.GetOutput(dataflow, config)
}

func GetDataflow(report types.Report, config settings.Config, isInternal bool) (*dataflow.DataFlow, *dataflow.DataFlow, error) {
	reportedDetections, _, err := detectors.GetOutput(report, config)
	if err != nil {
		return nil, nil, err
	}

	for _, detection := range reportedDetections {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}

	return dataflow.GetOutput(reportedDetections, config, isInternal)
}

func reportStats(report types.Report, config settings.Config) (*stats.Stats, *dataflow.DataFlow, error) {
	dataflowOutput, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, err
	}

	return stats.GetOutput(report.Inputgocloc, dataflowOutput, config)
}

func reportSecurity(
	report types.Report,
	config settings.Config,
	files []files.File,
	baseBranchFindings *basebranchfindings.Findings,
) (
	securityResults *security.Results,
	dataflow *dataflow.DataFlow,
	reportPassed bool,
	err error,
) {
	dataflow, _, err = GetDataflow(report, config, true)
	if err != nil {
		log.Debug().Msgf("error in dataflow %s", err)
		return
	}

	securityResults, reportPassed, err = security.GetOutput(dataflow, config, baseBranchFindings)
	if err != nil {
		log.Debug().Msgf("error in security %s", err)
		return
	}

	if config.Client != nil && config.Client.Error == nil {
		saas.SendReport(config, securityResults, files, dataflow)
	}

	return
}
