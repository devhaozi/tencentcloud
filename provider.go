package tencentcloud

import (
	"context"
	"errors"
	"strconv"

	"github.com/libdns/libdns"
)

func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	return p.listRecords(ctx, zone)
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for _, record := range records {
		if err := p.createRecord(ctx, zone, record); err != nil {
			return nil, err
		}
	}

	return records, nil
}

func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for _, record := range records {
		id, err := p.findRecord(ctx, zone, record)
		if err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				if err = p.createRecord(ctx, zone, record); err != nil {
					return nil, err
				}
				continue
			}
		}
		record.ID = strconv.FormatUint(id, 10)
		if err = p.modifyRecord(ctx, zone, record); err != nil {
			return nil, err
		}
	}

	return records, nil
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for _, record := range records {
		if record.ID == "" {
			id, err := p.findRecord(ctx, zone, record)
			if err != nil {
				return nil, err
			}
			record.ID = strconv.FormatUint(id, 10)
		}
		if err := p.deleteRecord(ctx, zone, record); err != nil {
			return nil, err
		}
	}

	return records, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
