package packer

func (opt option) Format() (verb rune) { return opt.verb }

func (opt option) Width() (width uint) { return opt.width }

func (opt option) Align() (align uint) { return opt.align }
