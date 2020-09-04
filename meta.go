package admin_field_duration

import (
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/ecletus/admin"
	"github.com/ecletus/core"
	"github.com/ecletus/core/resource"
	"github.com/hako/durafmt"
	"github.com/moisespsena-go/i18n-modular/i18nmod"
	path_helpers "github.com/moisespsena-go/path-helpers"
)

var (
	pkg      = path_helpers.GetCalledDir()
	i18ng    = i18nmod.PkgToGroup(pkg)
	unitsKey = "^" + i18ng + ".units"
)

func init() {
	admin.RegisterMetaConfigor("duration_string", func(meta *admin.Meta) {
		if meta.Config == nil {
			cfg := &DurationConfig{LimitN: 1}
			meta.Config = cfg
			cfg.ConfigureQorMeta(meta)
		}
	})

	admin.RegisterMetaTypeConfigor(reflect.TypeOf(time.Second), func(meta *admin.Meta) {
		if meta.Config == nil {
			cfg := &DurationConfig{LimitN: 1}
			meta.Config = cfg
			cfg.ConfigureQorMeta(meta)
		}
	})
}

type DurationConfig struct {
	// Non-zero to limit only first N elements to output.
	LimitN int
}

// ConfigureQorMeta configure select one meta
func (this *DurationConfig) ConfigureQorMeta(metaor resource.Metaor) {
	meta := metaor.(*admin.Meta)
	meta.Type = "duration_string"
	meta.ReadOnly = true
	meta.SetFormattedValuer(func(recorde interface{}, ctx *core.Context) interface{} {
		value := meta.Value(ctx, recorde)
		if value == nil {
			return nil
		}
		dur := value.(time.Duration)
		if dur == 0 {
			return ""
		}

		dura := durafmt.Parse(dur).LimitFirstN(this.LimitN)

		if unitsS := ctx.Ts(unitsKey); unitsS != unitsKey {
			if units, err := durafmt.DefaultUnitsCoder.Decode(unitsS); err != nil {
				ctx.AddError(errors.Wrapf(err, "parse %q as units", unitsS))
			} else {
				return dura.Format(units)
			}
		}
		return dura.String()
	})
}
