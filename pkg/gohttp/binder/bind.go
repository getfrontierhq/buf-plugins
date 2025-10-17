package binder

func (d *RequestDecoder) Bind(v interface{}) error {
	hasBody := hasBody(d.Request)

	d.BindHeader()
	if err := d.BindParams(v); err != nil {
		return err
	}
	if err := d.BindQuery(v); err != nil {
		return err
	}

	if !hasBody {
		return nil
	}

	if err := d.BindBody(v); err != nil {
		return err
	}

	return nil
}

func (e *RequestEncoder) Bind(v interface{}) error {
	shouldHaveBody := shouldHaveBody(e.Request.Method)

	e.BindHeader()
	if err := e.BindParams(v); err != nil {
		return err
	}
	if !shouldHaveBody {
		if err := e.BindQuery(v); err != nil {
			return err
		}
	} else {
		if err := e.BindBody(v); err != nil {
			return err
		}
	}

	return nil
}
