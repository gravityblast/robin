package robin

import "fmt"

type RunnerOptions struct {
	ExportersSize int
	LogLevel      int
}

var DefaultRunnerOptions = &RunnerOptions{
	ExportersSize: 1,
	LogLevel:      LogLevelFatal,
}

type scrapeJob struct {
	name string
}

type Runner struct {
	Scrapers    map[string]*Scraper
	options     *RunnerOptions
	exportQueue *exportQueue
	log         *applogger
}

func NewRunner(scrapers map[string]*Scraper, opts *RunnerOptions) *Runner {
	r := &Runner{
		Scrapers:    scrapers,
		options:     opts,
		exportQueue: newExportQueue(),
		log:         newAppLogger(opts.LogLevel),
	}

	return r
}

func (r *Runner) Run(name string, exp Exporter) error {
	var (
		s     *Scraper
		found bool
	)

	if s, found = r.Scrapers[name]; !found {
		return fmt.Errorf("unknown scraper `%s`", name)
	}

	r.log.Debug(fmt.Sprintf("run %d export workers", r.options.ExportersSize))
	r.exportQueue.run(r.options.ExportersSize)

	r.log.Debug(fmt.Sprintf("run scraper `%s`", name))
	err := s.Scrape(r.log, exp, r.exportQueue)
	if err != nil {
		return err
	}

	r.exportQueue.Wait()
	r.log.Debug(fmt.Sprintf("exporting done for scraper `%s`", name))

	return nil
}
